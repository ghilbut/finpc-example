import { sleep, check } from 'k6';
import grpc from 'k6/net/grpc';

export const options = {
    vus: 2,
    duration: '3s',
};

const client = new grpc.Client();
client.load([], '../grpc-server/.protos/v1/trading.proto');

export default async () => {
    client.connect('localhost:50051', {
        plaintext: true,
        timeout: '1s'
    });

    const response = client.invoke('trading.Trading/GetStockList', {})
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
};
