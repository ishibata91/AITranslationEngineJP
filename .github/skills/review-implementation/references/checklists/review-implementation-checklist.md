# Review Implementation Checklist

## Knowledge Check

- [ ] diff が owned_scope に収まっているか確認した
- [ ] product test と validation result を確認した
- [ ] backend を含む場合は mcp_mcp_docker_mcp-exec で python3 scripts/harness/run.py --suite all を実行し repo-local gate (coverage >= 70%, security=0, reliability=0, maintainability HIGH/BLOCKER=0) を確認した

## Common Pitfalls

- [ ] design review をしなかった
- [ ] 好みや将来改善で reroute しなかった
- [ ] 修正を同時に行わなかった
