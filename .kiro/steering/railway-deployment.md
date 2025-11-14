---
inclusion: always
---

# Railway Deployment for IOTA SDK

## Environment Configuration

### IOTA SDK Project
- **Project**: iota-sdk
- **Environments**: staging, production
- **Services**: api, database, redis

### Staging Database
```bash
Host: shuttle.proxy.rlwy.net
Port: 31150
Database: railway
User: postgres
Password: A6E4g1d2ae43Bebg2F65CEc3e56aa25g
```

## Common Operations

### Deployment
```bash
# Deploy to staging
railway up -s api -e staging --detach

# Redeploy latest
railway redeploy -s api -y

# Rollback
railway down -y
```

### Monitoring
```bash
# Tail logs
railway logs -s api --deployment

# SSH into container
railway ssh -s api -e staging

# Check status
railway status
```

### Database
```bash
# Connect to staging DB
railway connect database -e staging

# Direct connection
PGPASSWORD=A6E4g1d2ae43Bebg2F65CEc3e56aa25g \
  psql -h shuttle.proxy.rlwy.net -U postgres -p 31150 -d railway

# Run migrations
railway run -s api -e staging make db migrate up
```

## Safety Checklist

### Before Deployment
- Verify environment: `railway environment`
- Check service: `railway status`
- Review variables: `railway variables -s <service>`
- Run tests: `go test ./...`
- Backup database if production

### During Deployment
- Use `--detach` for CI/CD
- Monitor logs: `railway logs --deployment`
- Check health endpoints
- Verify service connectivity

### After Deployment
- Confirm success
- Check logs for errors
- Verify migrations applied
- Test API endpoints
- Monitor for 5 minutes
