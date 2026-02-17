# Run after: gh auth refresh -h github.com -s workflow
# (Visit the URL and complete the device flow)

$content = @'
name: CI
on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"
      - run: npm install
      - run: npm run lint
      - run: npm test
      - run: npm run build
'@
$b64 = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($content))
$body = @{message="chore: initialize TypeScript project structure"; content=$b64} | ConvertTo-Json
echo $body | gh api -X PUT "/repos/Clawland-AI/moltclaw/contents/.github/workflows/ci.yml" --input -
