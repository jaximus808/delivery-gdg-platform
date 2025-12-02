// lib/grpc-client.ts
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

const PROTO_PATH = path.join(process.cwd(), 'proto', 'order_service.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const orderProto = (grpc.loadPackageDefinition(packageDefinition) as any).order_service;


// Create gRPC client
const GRPC_SERVER_URL = process.env.GRPC_SERVER_URL || 'localhost:50051';

export const getOrderClient = () => {
  return new orderProto.OrderHandler(
    GRPC_SERVER_URL,
    grpc.credentials.createInsecure() // Use createSsl() for production with certificates
  );
};