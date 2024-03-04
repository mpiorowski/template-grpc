import protoLoader from "@grpc/proto-loader";
import { credentials, loadPackageDefinition } from "@grpc/grpc-js";
import { TARGET, AUTH_URI, PROFILE_URI } from "$env/static/private";

export const packageDefinition = protoLoader.loadSync(
    "./src/lib/proto/main.proto",
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
    },
);

const proto = /** @type {import("$lib/proto/main").ProtoGrpcType} */ (
    /** @type {unknown} */ (loadPackageDefinition(packageDefinition))
);

/** @type {import("@grpc/grpc-js").ChannelCredentials} */
const cr =
    TARGET === "production"
        ? credentials.createSsl()
        : credentials.createInsecure();

export const authService = new proto.proto.AuthService(AUTH_URI, cr);
export const profileService = new proto.proto.ProfileService(PROFILE_URI, cr);
