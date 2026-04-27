# Codex Review Trust Boundary Checklist

- [ ] 認証、認可、tenant isolation を確認した
- [ ] user-controlled input と外部 URL を確認した
- [ ] secret、admin 権限、PII を確認した
- [ ] SQL injection、XSS、SSRF、file upload を確認した
- [ ] violated invariant と root cause hypothesis を分けた
- [ ] local patch assessment と invariant tests を返した
- [ ] hard gate failure を他観点で相殺しなかった
