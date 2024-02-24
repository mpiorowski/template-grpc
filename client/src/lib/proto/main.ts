import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';

import type { UsersServiceClient as _proto_UsersServiceClient, UsersServiceDefinition as _proto_UsersServiceDefinition } from './proto/UsersService';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  proto: {
    AuthResponse: MessageTypeDefinition
    Count: MessageTypeDefinition
    Empty: MessageTypeDefinition
    Id: MessageTypeDefinition
    Page: MessageTypeDefinition
    StripeUrlResponse: MessageTypeDefinition
    User: MessageTypeDefinition
    UserRole: EnumTypeDefinition
    UsersService: SubtypeConstructor<typeof grpc.Client, _proto_UsersServiceClient> & { service: _proto_UsersServiceDefinition }
  }
}

