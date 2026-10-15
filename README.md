# pipeline-demo

Мини-проект на Go, показывающий полный конвейер доставки кода:
локальные хуки → CI → агент-ревьюер. Та же схема масштабируется на любой
реальный проект. Конвейер работает на GitHub Actions (`.github/workflows/`)
и на GitLab CI (`.gitlab-ci.yml`) — один и тот же набор проверок.

## Слои проверок

| Слой | Инструмент | Когда работает | Время |
|---|---|---|---|
| Локальные хуки | lefthook: go vet, golangci-lint, go test | перед каждым push | секунды |
| Линт | golangci-lint (staticcheck, revive, errcheck) | CI, stage `lint` | ~мин |
| Тесты | go test -race | CI, stage `test` | ~мин |
| Безопасность | gosec (SAST), gitleaks (секреты) | CI, stage `security` | ~мин |
| Ревью-агент | LLM по diff'у MR | вручную на MR | ~мин |

Принцип: детерминированные проверки (линтер, тесты, сканеры) — фундамент;
агент — слой суждения сверху. Агент только комментирует, никогда не мержит.

## Локальный запуск

    brew install lefthook
    lefthook install   # один раз в клоне репозитория

Теперь `git push` не пройдёт, пока не зелёные vet, lint и тесты.
Обход в экстренном случае: `git push --no-verify` — но CI всё равно всё поймает.

## Ревью-агент на PR/MR

Джоба `review-agent` берёт diff, отправляет в OpenAI-совместимый API и
постит разбор. На GitHub (Settings → Secrets and variables → Actions):

- секрет `LLM_API_KEY` — ключ API (без неё агент спит);
- Variable `LLM_BASE_URL` — например `https://api.openai.com/v1`;
- Variable `LLM_MODEL` — например `gpt-4o-mini`.

На GitLab — те же переменные через Settings → CI/CD → Variables, плюс
`MR_COMMENT_TOKEN` (project access token, scope `api`) для комментария в MR.

## Что демонстрирует «плохая» ветка

Ветка `feature/insecure-auth` содержит типовые проблемы: захардкоженный
пароль, `math/rand` для токенов, конкатенация SQL, непроверенная ошибка,
падающий тест. Локальный хук ловит это до push; если обойти (`--no-verify`)
— ловят golangci-lint, go test и gosec в CI.
