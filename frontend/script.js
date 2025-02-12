const API_URL = "http://127.0.0.1:8000/tasks"; // URL da API FastAPI

// Função para carregar tarefas da API
async function loadTasks() {
    const response = await fetch(API_URL);
    const tasks = await response.json();

    const taskList = document.getElementById("taskList");
    taskList.innerHTML = "";

    tasks.forEach(task => {
        const li = document.createElement("li");
        li.classList.add("task-item");

        // Criar o texto da tarefa com título e descrição
        const taskText = document.createElement("span");
        taskText.textContent = `${task.title} - ${task.description}`;

        // Criar botão de deletar
        const deleteButton = document.createElement("button");
        deleteButton.textContent = "🗑";
        deleteButton.classList.add("delete-button");
        deleteButton.onclick = () => deleteTask(task.id);

        // Adicionar elementos à lista
        li.appendChild(taskText);
        li.appendChild(deleteButton);
        taskList.appendChild(li);
    });
}

// Função para adicionar uma tarefa
async function addTask() {
    const title = document.getElementById("taskInput").value;
    const description = document.getElementById("taskDescription").value;

    if (!title || !description) {
        alert("Título e descrição são obrigatórios!");
        return;
    }

    await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, description }),
        mode: "cors"
    });

    document.getElementById("taskInput").value = "";
    document.getElementById("taskDescription").value = "";
    loadTasks();
}

// Função para excluir uma tarefa
async function deleteTask(id) {
    if (!confirm("Tem certeza que deseja excluir esta tarefa?")) return;
    await fetch(`${API_URL}/${id}`, { method: "DELETE" });
    loadTasks();
}

// Carregar tarefas ao abrir a página
loadTasks();