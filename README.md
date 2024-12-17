# i4
Industry 4.0 sample project

## Description of the repository
The project goal is to simulate a system which is 4.0 conformal. Client generated data is uploaded to the server and then saved for daily use and monitoring.

### Client
The client is written in Go and is made by 3 components:
- 2 data generators;
- a single data aggregator.
Each component can be customized using their configuration files.

The data generators simulate real working machines. Their data is saved locally on a simple text file.

The data aggregator simulate a PC living in the client network which periodically pulls the data from the machines using SSH connection, converts the data in the correct format and the sends it to the server.

In order to run the entire client using Docker just use the following command in the client folder (machine-client):
```bash
docker compose up
```
This will bring up all the components in the correct order and set them up as needed.

### Server
The server is written in Go and uses both Mongo and Kafka as support software. The server is made by 4(+1) components:
- a data collector;
- a data analyzer;
- Mongo DB;
- Kafka (and its configurator).
Each component can be customized using their configuration files.

The data collector is the server component which receives data from the client. After getting the file content, it submits an elaboration request on queue.

The data analyzer waits for elaboration requests. After getting one, it loads the file content from DB, splits it into single rows and saves them on the DB which restricted access for the given client.

In order to run the entire backend using Docker just use the following command in the backend folder (server):
```bash
docker compose up
```
This will bring up all the components in the correct order and set them up as needed.

## Next steps
The project next features could be:
- Add a way to visualize data from client (Grafana, for example);
- Add API keys for the client aggregators to use while communicating with the server.