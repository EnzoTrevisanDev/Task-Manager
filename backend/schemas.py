from pydantic import BaseModel

class Taskbase(BaseModel):
    title: str
    description: str
    
class TaskCreate(Taskbase):
    pass

class TaskResponse(Taskbase):
    id: int
    class Config:
        from_attributes = True
    