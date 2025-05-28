## **Twitter-Api**

API Backend de plataforma de microblogging similar a Twitter, implementado en Golang, con gestion de usuarios, manejo de tweets y suscripcion a feeds.


![Diagrama de arquitectura](/docs/diagram.png)


---

## **Levantarlo localmente**
   * Teniendo `docker` y `docker compose` instalado:
```bash
   docker compose up --build
```
---

# Endpoints disponibles
 
* Crear un usuario:
 
  **POST /users**
  ```json
  {
    "name": "...",
    "email": "..."
  }
  ```

  
* Usuario `:id` sigua a `:followID`:

  **POST /users/:id/follow/:followID**


* Usuario `:id` deje de seguir a `:followID`:

  **POST /users/:id/unfollow/:followID**


* Publicar un tweet:

  **POST /tweets**
    ```json
  {
    "user_id": 1,
    "content": "..."
  }
  ```

* Eliminar tweet:

  **DELETE /tweets/:id**


* Obtener feed:

  **GET /feed/:id**

