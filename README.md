# KanBan Card Service - MimoGen

Card management microservice for KanBan SaaS platform.

## Features

- Card CRUD with priority and due dates
- Label management
- Comment system
- Card assignment
- Move and reorder cards

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/cards | Create card |
| GET | /api/v1/cards/:id | Get card |
| PUT | /api/v1/cards/:id | Update card |
| POST | /api/v1/cards/:id/move | Move card |
| DELETE | /api/v1/cards/:id | Delete card |
| POST | /api/v1/labels | Create label |
| POST | /api/v1/comments | Create comment |

## Tech Stack

- Go 1.22
- Chi Router
- PostgreSQL

---

<details>
<summary><strong>IMPORTANT NOTICE</strong></summary>

<br>

**This repository was entirely generated using [MiMoCode](https://github.com/xiaomi/mimocode) - an AI-powered coding assistant by Xiaomi.**

All code, tests, documentation, and infrastructure configuration in this repository were created through AI-assisted development with MiMoCode.

---

**Этот репозиторий был полностью сгенерирован с помощью [MiMoCode](https://github.com/xiaomi/mimocode) - AI-ассистента для программирования от Xiaomi.**

Весь код, тесты, документация и инфраструктурная конфигурация в этом репозитории были созданы с помощью AI-ассистированной разработки MiMoCode.

</details>
