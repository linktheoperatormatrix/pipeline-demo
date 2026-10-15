#!/bin/sh
# Агент-ревьюер: берёт diff pull request'а (GitHub) или merge request'а
# (GitLab), отправляет в LLM и постит разбор комментарием (или в лог).
#
# GitHub (Settings → Secrets and variables → Actions):
#   LLM_API_KEY  — ключ OpenAI-совместимого API (Secret, включает агента)
#   LLM_BASE_URL — Variable, по умолчанию https://api.openai.com/v1
#   LLM_MODEL    — Variable, по умолчанию gpt-4o-mini
# GitLab (Settings → CI/CD → Variables):
#   то же + MR_COMMENT_TOKEN — project access token (scope api)
set -eu

if [ -z "${LLM_API_KEY:-}" ]; then
  echo '>> LLM_API_KEY не задан — агент спит.'
  echo '>> Добавь ключ в секреты репозитория и запусти джобу снова.'
  exit 0
fi

PLATFORM=none
if [ -n "${GITHUB_REPOSITORY:-}" ] && [ -n "${GH_TOKEN:-}" ]; then
  PLATFORM=github
elif [ -n "${CI_MERGE_REQUEST_IID:-}" ]; then
  PLATFORM=gitlab
fi

case "$PLATFORM" in
  github)
    PR=$(printf '%s' "${GITHUB_REF_NAME:-}" | cut -d/ -f1)
    case "$PR" in
      ''|*[!0-9]*) echo '>> Не pull request — пропуск.'; exit 0 ;;
    esac
    echo ">> PR #$PR: собираю diff..."
    DIFF=$(gh pr diff --repo "$GITHUB_REPOSITORY" "$PR" | head -c 50000)
    ;;
  gitlab)
    TARGET="origin/${CI_MERGE_REQUEST_TARGET_BRANCH_NAME:-main}"
    git fetch origin "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME" >/dev/null 2>&1 || true
    DIFF=$({ git diff --stat "$TARGET...HEAD"; echo
             git diff "$TARGET...HEAD"; } | head -c 50000)
    ;;
  *)
    echo '>> Нет PR/MR контекста — пропуск.'
    exit 0
    ;;
esac

PROMPT='Ты — опытный ревьюер Go-кода. Разбери diff и найди:
1. баги логики, гонки, утечки ресурсов;
2. безопасность: SQL-инъекции, слабая криптография, захардкоженные креды;
3. тесты: что не покрыто.
Отвечай по-русски, кратко, списком: [серъёзность] файл:строка — проблема и как исправить.
Если замечаний нет — напиши ровно: "Замечаний нет".'

echo '>> Отправляю diff в модель...'
REQUEST=$(jq -n \
  --arg model "${LLM_MODEL:-gpt-4o-mini}" \
  --arg prompt "$PROMPT" \
  --arg diff "$DIFF" \
  '{model: $model, messages: [{role: "system", content: $prompt},
                              {role: "user", content: ("Разбери этот diff:\n" + $diff)}],
    max_tokens: 1200, temperature: 0.2}')

RESPONSE=$(curl -sS --max-time 120 "${LLM_BASE_URL:-https://api.openai.com/v1}/chat/completions" \
  -H "Authorization: Bearer $LLM_API_KEY" \
  -H "Content-Type: application/json" \
  -d "$REQUEST")

REVIEW=$(printf '%s' "$RESPONSE" | jq -r \
  '.choices[0].message.content // .error.message // "модель вернула пустой ответ"')

echo
echo '================ РЕВЬЮ АГЕНТА ================'
printf '%s\n' "$REVIEW"
echo '==============================================='

case "$PLATFORM" in
  github)
    if printf '## 🤖 Ревью агента\n\n%s\n' "$REVIEW" |
       gh pr comment --repo "$GITHUB_REPOSITORY" "$PR" --body-file - >/dev/null; then
      echo '>> Ревью запощено комментарием в PR.'
    else
      echo '>> Не удалось запостить комментарий — ревью в логе джобы.'
    fi
    ;;
  gitlab)
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
    ;;
esac
