# Amazon-VL

[![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)](https://golang.org)
[![Tests](https://img.shields.io/badge/tests-7%20passed-success)](.)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Servidor HTTP leve para exposição segura de arquivos via web com autenticação HTTP Basic.

```
┌──────────────┐      ┌─────────────────────────────────┐      ┌────────────┐
│    Client    │─────▶│  amazon-vl (auth + fileserver)  │─────▶│   Files    │
│  curl/browser│◀─────│        :9000/healthz            │◀─────│  /var/log  │
└──────────────┘      └─────────────────────────────────┘      └────────────┘
```

## Quick Start

```bash
# Build
make build

# Run
./bin/amazon-vl /path/to/logs 9000

# Access
curl -u joaquim:amazon http://localhost:9000/
```

## Features

| Feature | Descrição |
|---------|-----------|
| 🔐 **Auth** | HTTP Basic com MD5 crypt hash |
| 🏥 **Health** | Endpoint `/healthz` para Kubernetes |
| 🛑 **Graceful** | Shutdown limpo via SIGTERM |
| 📊 **Logging** | Access logs estruturados |
| ⚡ **Timeouts** | Read/Write/Idle configurados |
| 🐳 **Docker** | Multi-stage build pronto |

## Instalação

### Build Local

```bash
git clone https://github.com/joaquimsnjunior/amazon-vl.git
cd amazon-vl
make build
```

### Docker

```bash
make docker-build
docker run -d -p 9000:9000 -v /var/log:/logs:ro amazon-vl:latest /logs 9000
```

## Uso

```bash
amazon-vl [OPTIONS] <directory> <port>

ARGUMENTS:
    <directory>    Diretório a ser servido
    <port>         Porta HTTP (ex: 8080, 9000)

OPTIONS:
    --help         Mostra ajuda
    --version      Mostra versão
```

### Exemplos

```bash
# Servir /var/log na porta 9000
./bin/amazon-vl /var/log 9000

# Com credenciais customizadas
AUTH_USER=admin AUTH_HASH='$1$xyz...' ./bin/amazon-vl /var/log 9000

# Verificar health
curl http://localhost:9000/healthz
```

## Configuração

### Variáveis de Ambiente

| Variável | Default | Descrição |
|----------|---------|-----------|
| `AUTH_USER` | `joaquim` | Usuário para autenticação |
| `AUTH_HASH` | `$1$neD...` | Hash MD5 crypt da senha |
| `AUTH_REALM` | `amazon-server-logs.com` | Realm do Basic Auth |

### Gerar Hash de Senha

```bash
# Via script incluído
./scripts/generate-hash.sh minhasenha

# Via openssl
openssl passwd -1 -salt "$(openssl rand -base64 6)" "minhasenha"
```

## Estrutura do Projeto

```
amazon-vl/
├── cmd/
│   └── main.go              # Entrypoint
├── internal/
│   ├── auth/
│   │   ├── basic.go         # Autenticação
│   │   └── basic_test.go
│   └── server/
│       ├── handler.go       # FileServer handler
│       ├── server.go        # HTTP server + graceful shutdown
│       └── server_test.go
├── configs/
│   ├── .env.example
│   └── .htpasswd.example
├── scripts/
│   └── generate-hash.sh
├── Dockerfile
├── Makefile
└── go.mod
```

## Desenvolvimento

```bash
# Instalar dependências
go mod tidy

# Rodar testes
make test

# Rodar com coverage
make test-coverage

# Lint
make lint

# Formatar código
make fmt
```

## Deploy

### Systemd

```ini
# /etc/systemd/system/amazon-vl.service
[Unit]
Description=Amazon Log Viewer
After=network.target

[Service]
Type=simple
User=logviewer
ExecStart=/usr/local/bin/amazon-vl /var/log/app 9000
Restart=on-failure
Environment=AUTH_USER=admin
Environment=AUTH_HASH=$1$...

[Install]
WantedBy=multi-user.target
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: amazon-vl
spec:
  template:
    spec:
      containers:
      - name: amazon-vl
        image: amazon-vl:1.1.0
        args: ["/logs", "9000"]
        ports:
        - containerPort: 9000
        env:
        - name: AUTH_USER
          valueFrom:
            secretKeyRef:
              name: amazon-vl-auth
              key: username
        - name: AUTH_HASH
          valueFrom:
            secretKeyRef:
              name: amazon-vl-auth
              key: hash
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9000
          initialDelaySeconds: 5
        volumeMounts:
        - name: logs
          mountPath: /logs
          readOnly: true
```

## Makefile

```bash
make help           # Ver comandos disponíveis
make build          # Compilar binário
make build-static   # Compilar binário estático (containers)
make run            # Executar (requer DIR e PORT)
make test           # Rodar testes
make test-coverage  # Testes com coverage
make docker-build   # Build imagem Docker
make docker-run     # Rodar container
make clean          # Limpar artefatos
make install        # Instalar em /usr/local/bin
```

## Segurança

- ✅ Credenciais externalizadas via env vars
- ✅ Container roda como non-root (UID 1000)
- ✅ Suporte a volume read-only
- ✅ Timeouts HTTP configurados
- ⚠️ Recomenda-se TLS via reverse proxy (nginx/traefik)

### Produção Recomendada

```nginx
server {
    listen 443 ssl;
    server_name logs.example.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://127.0.0.1:9000;
    }
}
```

## API

| Endpoint | Auth | Descrição |
|----------|------|-----------|
| `GET /` | ✅ | Lista arquivos do diretório |
| `GET /{path}` | ✅ | Serve arquivo/diretório |
| `GET /healthz` | ❌ | Health check (retorna `{"status":"healthy"}`) |

## Troubleshooting

| Problema | Causa | Solução |
|----------|-------|---------|
| `address already in use` | Porta ocupada | `lsof -i :9000` e matar processo |
| `401 Unauthorized` | Credenciais erradas | Verificar AUTH_USER/AUTH_HASH |
| `permission denied` | Sem acesso ao dir | Verificar permissões do usuário |

## Contribuição

```bash
# Fork e clone
git clone https://github.com/your-user/amazon-vl.git

# Criar branch
git checkout -b feature/nova-feature

# Desenvolver, testar, commitar
make test
git commit -m "feat: adiciona feature X"

# Push e PR
git push origin feature/nova-feature
```

## Licença

MIT License - veja [LICENSE](LICENSE)

---

<p align="center">
  <sub>Built with ☕ for SREs</sub>
</p>
