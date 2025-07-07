# CRUD API RESTful em Go (Golang)

Este projeto é uma API RESTful escrita em Go (Golang), utilizando o framework [Chi](https://github.com/go-chi/chi). A aplicação implementa operações básicas de CRUD (Create, Read, Update, Delete) para gerenciamento de usuários em memória.

---

## 🚀 Funcionalidades

- ✅ Criar um novo usuário (POST)
- 📖 Listar todos os usuários (GET)
- 🔍 Obter um usuário por ID (GET)
- ✏️ Atualizar um usuário existente (PUT)
- ❌ Deletar um usuário (DELETE)

---

## 🛠️ Tecnologias Utilizadas

- [Go (Golang)](https://golang.org/)
- [Chi Router](https://github.com/go-chi/chi)
- [UUID](https://github.com/google/uuid)
- JSON nativo para serialização

---

📬 Endpoints da API
| Método | Endpoint          | Descrição                     |
| ------ | ----------------- | ----------------------------- |
| POST   | `/api/users`      | Cria um novo usuário          |
| GET    | `/api/users`      | Retorna todos os usuários     |
| GET    | `/api/users/{id}` | Retorna um usuário por ID     |
| PUT    | `/api/users/{id}` | Atualiza um usuário existente |
| DELETE | `/api/users/{id}` | Remove um usuário por ID      |

---

🧠 Observações
Os dados dos usuários são armazenados em memória (não persistem entre execuções).

O campo id é gerado automaticamente via UUID.
