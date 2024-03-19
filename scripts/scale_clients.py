# Script to scale the number of client containers in docker-compose file
import sys
import yaml

DOCKER_COMPOSE_FILE_PATH = "../docker-compose-dev.yaml"
CONFIG_PARAM_CONTAINER_NAME = "container_name"
CONFIG_PARAM_ENVIRONMENT = "environment"
CLIENT_SERVICE_PREFIX = "client"
CLIENT_SERVICE_ENV_ID_NAME = "CLIENT_ID"
FILE_BEGIN_POSITION = 0

def get_service_client_config():
    return {
        CONFIG_PARAM_CONTAINER_NAME: "",
        "image": "client:latest",
        "entrypoint": "/client",
        CONFIG_PARAM_ENVIRONMENT: [
            "CLI_LOG_LEVEL=DEBUG"
        ],
        "networks": [
            "testing_net"
        ],
        "depends_on": [
            "server"
        ]
    }

def main():
    try:
        _, new_number_of_client_containers = sys.argv
        new_number_of_client_containers = int(new_number_of_client_containers)
        with open(DOCKER_COMPOSE_FILE_PATH, "r+") as docker_compose_file:
            docker_compose_file_as_yaml = yaml.safe_load(docker_compose_file)
            docker_services = docker_compose_file_as_yaml["services"]
            client_containers = [
                (service, config) for service, config in docker_services.items() if service.startswith(CLIENT_SERVICE_PREFIX)
            ]
            if (new_number_of_client_containers > len(client_containers)):
                for i in range(len(client_containers), new_number_of_client_containers):
                    new_client_service_name = f"{CLIENT_SERVICE_PREFIX}_{i + 1}"
                    new_client_service_config = get_service_client_config()
                    new_client_service_config[CONFIG_PARAM_CONTAINER_NAME] = new_client_service_name
                    new_client_service_config[CONFIG_PARAM_ENVIRONMENT].append(f"{CLIENT_SERVICE_ENV_ID_NAME}={i + 1}")
                    docker_services[new_client_service_name] = new_client_service_config
            elif (new_number_of_client_containers < len(client_containers)):
                for i in range(len(client_containers), new_number_of_client_containers, -1):
                    client_service_name = f"{CLIENT_SERVICE_PREFIX}_{i}"
                    del docker_services[client_service_name]

            docker_compose_file.seek(FILE_BEGIN_POSITION)
            docker_compose_file.truncate()
            yaml.dump(docker_compose_file_as_yaml, docker_compose_file, sort_keys=False)
    except ValueError as error:
        script_name, *rest = sys.argv
        print(f"{error.__class__.__name__}: {error}")
        print(f"Usage: poetry run python {script_name} <number_of_client_containers>")
        sys.exit(1)
    except OSError as error:
        print(f"{error.__class__.__name__}: {error}")
        sys.exit(1)
    except Exception as error:
        print(f"Unexpected error occurred: {error}")
        sys.exit(1)

main()
