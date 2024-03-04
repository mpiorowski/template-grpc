export type Safe<T> =
    | {
          success: true;
          data: T;
      }
    | {
          success: false;
          error: string;
          fields?: {
              field: string;
              tag: string;
          }[];
      };

export declare function safe<T>(promise: Promise<T>): Promise<Safe<T>>;
export declare function safe<T>(fn: () => T): Safe<T>;
