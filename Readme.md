# Gerenciador de Tarefas

Um gerenciador de tarefas simples e eficiente, desenvolvido com **FastAPI (Python) no backend** e **HTML, CSS e JavaScript no frontend**. O objetivo deste projeto é demonstrar a implementação de um CRUD básico para gestão de tarefas.

## 🚀 Tecnologias Utilizadas
- **Backend:** Python + FastAPI
- **Frontend:** HTML + CSS + JavaScript
- **Banco de Dados:** SQLite
- **Hospedagem:** Pode ser rodado localmente ou implantado em uma VPS/Heroku/Railway

## 📌 Funcionalidades
✅ Criar tarefas com título e descrição  
✅ Listar tarefas  
✅ Deletar tarefas  
✅ Interface minimalista e responsiva  

---

## 📥 Como Rodar o Projeto

### 1️⃣ Clonar o Repositório
```sh
git clone https://github.com/seuusuario/gerenciador-tarefas.git
cd gerenciador-tarefas
```

### 2️⃣ Criar e Ativar o Ambiente Virtual
```sh
python -m venv venv  # Criar ambiente virtual
source venv/bin/activate  # Ativar no Linux/macOS
venv\Scripts\activate  # Ativar no Windows
```

### 3️⃣ Instalar as Dependências
```sh
pip install -r requirements.txt
```

### 4️⃣ Iniciar o Servidor FastAPI
```sh
uvicorn main:app --reload
```
A API estará rodando em: `http://127.0.0.1:8000`

### 5️⃣ Abrir o Frontend no Navegador
Basta abrir o arquivo `index.html` no navegador ou rodar via um servidor local.

---

## 🛠️ Endpoints da API
| Método | Endpoint          | Descrição |
|---------|-----------------|-------------|
| GET     | `/tasks`        | Lista todas as tarefas |
| POST    | `/tasks`        | Adiciona uma nova tarefa |
| DELETE  | `/tasks/{id}`   | Exclui uma tarefa |

---

## 📸 Demonstração
![Task Manager Preview](caminho-do-gif.gif)  
> GIF mostrando a aplicação funcionando (instruções abaixo para gerar)

---

## 📄 Licença
Este projeto está licenciado sob a **MIT License** - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## ✨ Autor
**Enzo Trevisan** - [LinkedIn](https://www.linkedin.com/in/seu-perfil/) - [GitHub](https://github.com/seuusuario/)

