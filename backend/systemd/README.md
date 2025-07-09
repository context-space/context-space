## Context Space Backend Service

Save this configuration:

```bash
cp context-space-backend.service /etc/systemd/system/context-space-backend.service
```

Run these commands:

```bash
sudo systemctl daemon-reload
sudo systemctl enable context-space-backend
sudo systemctl start context-space-backend
```

To check the service status:

```bash
sudo systemctl status context-space-backend
```
