# Tarea 2!!

## Consideraciones generales
- Antes de correr el programa, se deben setear las variables de entorno de la siguiente manera:
    ```yaml
    MV2_IP: "127.0.0.1"
    MV3_IP: "127.0.0.1"
    MONGO_PORT: 27017
    GRPC_PORT: 50051
    RABBIT_PORT: 5672
    ```
- Luego se deben copiar estas variables a cada uno de los servicios. Para agilizar, si se encuentra en Linux usar el script `copy_env.sh` (no olvidar dar permisos de ejecucion con `chmod +x copy_env.sh`). Si se encuentra en Windows, copiar las variables de entorno manualmente.

- Para correr los servicios debe acceder a las subcarpetas de cada uno de ellos para correrlos de manera individual.
    ```shell
    cd inventario
    go run inventario.go
    ```

- De ser necesario no olvidar correr `go mod download` para descargar las dependencias de cada servicio. 