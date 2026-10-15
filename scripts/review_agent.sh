#!/bin/sh
# Агент-ревьюер: берёт diff merge request'а, отправляет в LLM
# и постит разбор комментарием в MR (или печатает в лог джобы).
#
# Переменные CI/CD (Settings → CI/CD → Variables):
#   LLM_API_KEY      — ключ OpenAI-совместимого API (masked) — включает агента
#   LLM_BASE_URL     — например https://api.openai.com/v1
#   LLM_MODEL        — например gpt-4o-mini
#   MR_COMMENT_TOKEN — project access token (scope api): постить комментарий в MR
set -eu

if [ -z "${LLM_API_KEY:-}" ]; then
  echo '>> LLM_API_KEY не задан — агент спит.'
  echo '>> Добавь ключ в Settings → CI/CD → Variables и запусти джобу снова.'
  exit 0
fi
if [ -z "${CI_MERGE_REQUEST_IID:-}" ]; then
  echo '>> Не merge request пайплайн — пропуск.'
  exit 0
fi

TARGET="origin/${CI_MERGE_REQUEST_TARGET_BRANCH_NAME:-main}"
git fetch origin "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME" >/dev/null 2>&1 || true

# diff + статистика; режем до 50 КБ, чтобы влезло в контекст модели
DIFF=$({ git diff --stat "$TARGET...HEAD"; echo; git diff "$TARGET...HEAD"; } | head -c 50000)

PROMPT='Ты — опытный ревьюер Go-кода. Разбери diff merge request'а и найди:
1. баги логики, гонки, утечки ресурсов;
2. безопасность: SQL-инъекции, слабая криптография, захардкоженные креды;
3. тесты: что не покрыто.
Отвечай по-русски, кратко, списком: [серъёзность] файл:строка — проблема и как исправить.
Если замечаний нет — напиши ровно: "Замечаний нет".'

echo '>> Отправляю diff в модель...'
REQUEST=$(jq -n \
  --arg model "$LLM_MODEL" \
  --arg prompt "$PROMPT" \
  --arg diff "$DIFF" \
  '{model: $model, messages: [{role: "system", content: $prompt},
                              {role: "user", content: ("Разбери этот diff:\n" + $diff)}],
    max_tokens: 1200, temperature: 0.2}')

RESPONSE=$(curl -sS --max-time 120 "$LLM_BASE_URL/chat/completions" \
  -H "Authorization: Bearer $LLM_API_KEY" \
  -H "Content-Type: application/json" \
  -d "$REQUEST")

REVIEW=$(printf '%s' "$RESPONSE" | jq -r '.choices[0].message.content // .error.message // "модель вернула пустой ответ"')

echo
echo '================ РЕВЬЮ АГЕНТА ================'
printf '%s\n' "$REVIEW"
echo '==============================================='

if [ -n "${MR_COMMENT_TOKEN:-}" ]; then
  BODY=$(jq -n --arg body "## 🤖 Ревью агента

$REVIEW" '{body: $body}')
  curl -sS --max-time 60 --fail \
    --header "PRIVATE-TOKEN: $MR_COMMENT_TOKEN" \
    --data "$BODY" \
    "$CI_API_V4_URL/projects/$CI_PROJECT_ID/merge_requests/$CI_MERGE_REQUEST_IID/notes" >/dev/null
  echo '>> Ревью запощено комментарием в MR.'
else
  echo '>> MR_COMMENT_TOKEN не задан — ревью только в логе джобы.'
fi
