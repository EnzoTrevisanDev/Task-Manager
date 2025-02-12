from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy.orm import Session
from database import SessionLocal, engine, init_db
from models import Task
from schemas import TaskCreate, TaskResponse
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

# Habilitar CORS (permite que o frontend acesse a API)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Pode restringir para ["http://127.0.0.1:5500"] se estiver rodando local
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

init_db()

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


# Create a task
@app.post("/tasks", response_model=TaskResponse)
def create_task(task: TaskCreate, db: Session = Depends(get_db)):
    db_task = Task(title=task.title, description=task.description)
    db.add(db_task)
    db.commit()
    db.refresh(db_task)
    return db_task

#list all task
@app.get("/tasks", response_model=list[TaskResponse])
def read_tasks(db: Session = Depends(get_db)):
    return db.query(Task).all()

#delete a single task
@app.delete("/tasks/{task_id}")
def delete_task(task_id: int, db: Session = Depends(get_db)):
    db_task = db.query(Task).filter(Task.id == task_id).first()
    if db_task is None:
        raise HTTPException(status_code=404, detail="Tarefa não encontrada")
    db.delete(db_task)
    db.commit()
    return {"message": "Tarefa deletada com sucesso"}