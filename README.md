# ESP_MESH_TOPOLOGY

API desarrollada en Go que recibe paquetes de una red en malla de dispositivos ESP32 o ESP8266 y elabora una representación gráfica, en forma de grafo, que describe la topología de la red y exhibe la información de cada nodo.
Está pensada para nodos constituidos de ESP32 con sensores de temperatura y humedad.

## Único endpoint: POST /v1/update

- Recibe el reporte de un nodo de la red mesh (MAC del propio nodo y su padre, temperatura, humedad y si es el nodo raíz). El servidor acumula estos reportes en un buffer hasta reunir el número esperado de nodos y luego genera la imagen de la topología.

- URL base (por defecto): http://localhost:8080

### La petición en cuestión (heh)

Parámetros:
- `parentMAC` (string): MAC del nodo padre en la malla.
- `selfMAC` (string): MAC del propio nodo.
- `temp` (number, float): Temperatura.
- `humidity` (number, float): Humedad.
- `isRoot` (boolean): Indica si el nodo es el raíz.

Ejemplo:
```json
{
  "parentMAC": "AA:BB:CC:DD:EE:01",
  "selfMAC":   "AA:BB:CC:DD:EE:02",
  "temp":      24.7,
  "humidity":  52.3,
  "isRoot":    false
}
```

### Cómo usarlo (curl)

```bash
curl -X POST "http://localhost:8080/v1/update" \
  -H "Content-Type: application/json" \
  -d '{
        "parentMAC": "AA:BB:CC:DD:EE:01",
        "selfMAC":   "AA:BB:CC:DD:EE:02",
        "temp":      24.7,
        "humidity":  52.3,
        "isRoot":    false
      }'
```

### Respuestas

- 200 OK (espero por más paquetes):
  ```json
  {
    "status": "Nodo agregado a buffer, esperando por X paquetes más.",
    "nodes":  N
  }
  ```
  Donde `X` es la cantidad de paquetes restantes para completar el buffer y `N` el número de nodos acumulados.

- 200 OK (topología generada):
  ```json
  {
    "status": "graph updated",
    "nodes":  N
  }
  ```
  Al completar el buffer, se genera la imagen `static/topology.png` que puede consultarse vía `http://localhost:8080/static/topology.png` (la página `http://localhost:8080/static/` actualiza la imagen automáticamente cada 5 segundos).

- 400 Bad Request:
  ```json
  { "error": "uuuu" }
  ```

## Ejemplo de grafo

![Uuuuu](./static/topology.png)
