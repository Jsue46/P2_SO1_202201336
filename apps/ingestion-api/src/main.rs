use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::env;

#[derive(Debug, Deserialize, Serialize)]
struct SolicitudVenta {
    categoria: String,
    producto: String,
    precio: f64,
    cantidad_vendida: i32,
}

#[get("/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().body("OK")
}

#[post("/sale")]
async fn crear_venta(
    cuerpo_req: web::Json<SolicitudVenta>,
    cliente: web::Data<Client>,
) -> impl Responder {
    let url_go = env::var("GO_API_URL").unwrap_or_else(|_| "http://localhost:8080".to_string());

    // Asumimos que el servicio de Go expone un endpoint POST /sale
    let res = cliente
        .post(format!("{}/sale", url_go))
        .json(&cuerpo_req)
        .send()
        .await;

    match res {
        Ok(respuesta) => {
            if respuesta.status().is_success() {
                HttpResponse::Ok().json(serde_json::json!({"estado": "enviado a go"}))
            } else {
                println!("Servicio Go retornó error: {:?}", respuesta.status());
                HttpResponse::InternalServerError().body("Error en servicio upstream")
            }
        }
        Err(e) => {
            println!("Falló conexión con Servicio Go: {}", e);
            HttpResponse::InternalServerError().body("Fallo al conectar con Servicio Go")
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    let cliente = Client::new();

    let puerto_str = env::var("PORT").unwrap_or_else(|_| "5000".to_string());
    let puerto = puerto_str.parse::<u16>().unwrap_or(5000);

    println!("Iniciando API Rust en puerto {}", puerto);

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(cliente.clone()))
            .service(crear_venta)
            .service(health_check)
    })
    .bind(("0.0.0.0", puerto))?
    .run()
    .await
}
