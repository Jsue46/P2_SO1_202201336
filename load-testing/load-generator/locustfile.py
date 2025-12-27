from locust import HttpUser, task, between
import random

class UsuarioBlackFriday(HttpUser):
    wait_time = between(0.5, 2)

    @task
    def enviar_venta(self):
        categorias = ["Electronica", "Ropa", "Hogar", "Belleza"]
        productos = {
            "Electronica": ["Laptop", "Smartphone", "TV", "Audifonos", "Camara"],
            "Ropa": ["Camisa", "Pantalones", "Vestido", "Chaqueta", "Zapatos"],
            "Hogar": ["Mesa", "Silla", "Lampara", "Cama", "Sofa"],
            "Belleza": ["Labial", "Perfume", "Crema", "Shampoo", "Maquillaje"]
        }
        
        categoria = random.choice(categorias)
        producto = random.choice(productos[categoria])
        precio = round(random.uniform(10.0, 1000.0), 2)
        qty = random.randint(1, 5)
        
        payload = {
            "categoria": categoria,
            "producto": producto,
            "precio": precio,
            "cantidad_vendida": qty
        }
        
        # El endpoint debe coincidir con la ruta de Rust/Ingress
        # Ahora es /sale (actualizado en Rust API)
        self.client.post("/sale", json=payload)
