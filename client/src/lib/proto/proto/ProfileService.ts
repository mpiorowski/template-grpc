// Original file: proto/main.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { Empty as _proto_Empty, Empty__Output as _proto_Empty__Output } from '../proto/Empty';
import type { Profile as _proto_Profile, Profile__Output as _proto_Profile__Output } from '../proto/Profile';

export interface ProfileServiceClient extends grpc.Client {
  GetProfile(argument: _proto_Empty, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  GetProfile(argument: _proto_Empty, metadata: grpc.Metadata, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  GetProfile(argument: _proto_Empty, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  GetProfile(argument: _proto_Empty, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  getProfile(argument: _proto_Empty, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  getProfile(argument: _proto_Empty, metadata: grpc.Metadata, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  getProfile(argument: _proto_Empty, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  getProfile(argument: _proto_Empty, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  
  UpdateProfile(argument: _proto_Profile, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  UpdateProfile(argument: _proto_Profile, metadata: grpc.Metadata, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  UpdateProfile(argument: _proto_Profile, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  UpdateProfile(argument: _proto_Profile, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  updateProfile(argument: _proto_Profile, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  updateProfile(argument: _proto_Profile, metadata: grpc.Metadata, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  updateProfile(argument: _proto_Profile, options: grpc.CallOptions, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  updateProfile(argument: _proto_Profile, callback: grpc.requestCallback<_proto_Profile__Output>): grpc.ClientUnaryCall;
  
}

export interface ProfileServiceHandlers extends grpc.UntypedServiceImplementation {
  GetProfile: grpc.handleUnaryCall<_proto_Empty__Output, _proto_Profile>;
  
  UpdateProfile: grpc.handleUnaryCall<_proto_Profile__Output, _proto_Profile>;
  
}

export interface ProfileServiceDefinition extends grpc.ServiceDefinition {
  GetProfile: MethodDefinition<_proto_Empty, _proto_Profile, _proto_Empty__Output, _proto_Profile__Output>
  UpdateProfile: MethodDefinition<_proto_Profile, _proto_Profile, _proto_Profile__Output, _proto_Profile__Output>
}
