export interface paths {
    "/ai-events": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List AI events */
        get: {
            parameters: {
                query?: {
                    /** @description Filter by license plate */
                    license_plate?: string;
                    /** @description Filter by gate ID */
                    gate_id?: string;
                    /** @description Filter by event type */
                    event_type?: string;
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create an AI event */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description AI event payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.CreateAIEventRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.AIEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/ai-events/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get an AI event by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description AI event ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.AIEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete an AI event */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description AI event ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update an AI event */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description AI event ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description AI event update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.UpdateAIEventRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.AIEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/billing": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List billing records */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a billing record */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Billing payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.CreateBillingRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/billing/code/{billing_code}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a billing record by billing code */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Billing code */
                    billing_code: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/billing/order/{order_code}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a billing record by order code */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Order code */
                    order_code: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/billing/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a billing record by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Billing ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a billing record */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Billing ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a billing record */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Billing ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Billing update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.UpdateBillingRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/drivers": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List drivers */
        get: {
            parameters: {
                query?: {
                    /** @description Filter by driver status */
                    status?: string;
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a driver */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Driver payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.CreateDriverRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/drivers/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a driver by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Driver ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a driver */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Driver ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a driver */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Driver ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Driver update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.UpdateDriverRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/inbounds": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List inbound receipts */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create an inbound receipt */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Inbound payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.CreateInboundRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/inbounds/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get an inbound receipt by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inbound ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete an inbound receipt */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inbound ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update an inbound receipt */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inbound ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Inbound update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.UpdateInboundRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/inventory": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List inventory records */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create an inventory record */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Inventory payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.CreateInventoryRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/inventory/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get an inventory record by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inventory ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete an inventory record */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inventory ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update an inventory record */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Inventory ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Inventory update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.UpdateInventoryRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/orders": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List orders */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create an order */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Order payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.CreateOrderRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/orders/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get an order by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Order ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete an order */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Order ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update an order */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Order ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Order update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.UpdateOrderRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/products": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List products */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a product */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Product payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.CreateProductRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/products/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a product by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Product ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a product */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Product ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a product */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Product ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Product update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.UpdateProductRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/profile": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get current user profile */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        /** Update current user profile */
        put: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Profile update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.UpdateProfileRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        /** Create user profile */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Create user profile */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.CreateProfileRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/roles/me": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get roles and permissions for current user */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_role_dto.UserRolesAndPermissionsResponse"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/tracking-logs": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List tracking events */
        get: {
            parameters: {
                query?: {
                    /** @description Filter by order code */
                    order_code?: string;
                    /** @description Filter by driver code */
                    driver_code?: string;
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a tracking event */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Tracking event payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.CreateTrackingEventRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/tracking-logs/order/{order_code}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List tracking events for an order */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Order code */
                    order_code: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"][];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/tracking-logs/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a tracking event by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Tracking event ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        /** Update a tracking event */
        put: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Tracking event ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Tracking event update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.UpdateTrackingEventRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        post?: never;
        /** Delete a tracking event */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Tracking event ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/trips": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List trips */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a trip */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Trip payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.CreateTripRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/trips/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a trip by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Trip ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a trip */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Trip ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a trip */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Trip ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Trip update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.UpdateTripRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/trips/{id}/assign-driver": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Assign a driver to a trip */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Trip ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Driver code */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.AssignDriverRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List users */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a user */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description User payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.CreateUserRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a user by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description User ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a user */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description User ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a user */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description User ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description User update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UpdateUserRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/warehouses": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List warehouses */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a warehouse */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Warehouse payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.CreateWarehouseRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/warehouses/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a warehouse by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Warehouse ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a warehouse */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Warehouse ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a warehouse */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Warehouse ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Warehouse update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.UpdateWarehouseRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
    "/workspaces": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List workspaces */
        get: {
            parameters: {
                query?: {
                    /** @description Number of items to skip */
                    skip?: number;
                    /** @description Max items per page */
                    limit?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.PaginatedResponse"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a workspace */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Workspace payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.CreateWorkspaceRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.WorkspaceResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/workspaces/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a workspace by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Workspace ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.WorkspaceResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /** Delete a workspace */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Workspace ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description No Content */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "*/*": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /** Update a workspace */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Workspace ID */
                    id: number;
                };
                cookie?: never;
            };
            /** @description Workspace update */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.UpdateWorkspaceRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.WorkspaceResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            [key: string]: unknown;
                        };
                    };
                };
            };
        };
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        "my-web-app_com_smart-logistic-hub_internal_ai_dto.AIEventResponse": {
            confidence_score?: number;
            created_at?: string;
            event_code?: string;
            event_type?: string;
            gate_id?: string;
            id?: number;
            license_plate?: string;
            low_confidence?: boolean;
            matched_driver_id?: number;
            matched_trip_id?: number;
            timestamp?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_ai_dto.CreateAIEventRequest": {
            confidence_score?: number;
            event_code?: string;
            /** @enum {string} */
            event_type: "INBOUND" | "OUTBOUND";
            gate_id: string;
            license_plate: string;
            timestamp?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_ai_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_ai_dto.AIEventResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_ai_dto.UpdateAIEventRequest": {
            confidence_score?: number;
            event_type?: string;
            gate_id?: string;
            license_plate?: string;
            timestamp?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse": {
            amount_total?: number;
            billing_code?: string;
            created_at?: string;
            created_by?: string;
            currency?: string;
            id?: number;
            order_code?: string;
            paid_at?: string;
            payer_email?: string;
            payer_name?: string;
            payer_phone?: string;
            payment_method?: string;
            payment_status?: string;
            transaction_id?: string;
            updated_at?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_billing_dto.CreateBillingRequest": {
            amount_total?: number;
            billing_code: string;
            created_by?: string;
            currency?: string;
            order_code: string;
            payer_info: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.PayerInfo"];
            payment_method?: string;
            payment_status?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_billing_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_billing_dto.BillingResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_billing_dto.PayerInfo": {
            email?: string;
            name: string;
            phone: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_billing_dto.UpdateBillingRequest": {
            paid_at?: string;
            payment_status?: string;
            transaction_id?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_driver_dto.CreateDriverRequest": {
            current_lat?: number;
            current_lng?: number;
            driver_code: string;
            full_name: string;
            license_plate: string;
            phone: string;
            /** @enum {string} */
            status?: "AVAILABLE" | "BUSY" | "OFFLINE";
            vehicle_type: string;
            warehouse_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse": {
            created_at?: string;
            created_by?: string;
            current_lat?: number;
            current_lng?: number;
            driver_code?: string;
            full_name?: string;
            id?: number;
            license_plate?: string;
            phone?: string;
            status?: string;
            updated_at?: string;
            vehicle_type?: string;
            warehouse_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_driver_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_driver_dto.UpdateDriverRequest": {
            current_lat?: number;
            current_lng?: number;
            driver_code?: string;
            full_name?: string;
            license_plate?: string;
            phone?: string;
            /** @enum {string} */
            status?: "AVAILABLE" | "BUSY" | "OFFLINE";
            vehicle_type?: string;
            warehouse_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.CreateInboundRequest": {
            created_by?: string;
            items: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemRequest"][];
            receipt_code: string;
            status?: string;
            supplier_name: string;
            warehouse_id: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemRequest": {
            expected_qty?: number;
            product_id: number;
            qc_passed?: number;
            received_qty?: number;
            rejected_qty?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemResponse": {
            expected_qty?: number;
            id?: number;
            inbound_id?: number;
            product_id?: number;
            qc_passed?: number;
            received_qty?: number;
            rejected_qty?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse": {
            completed_at?: string;
            created_at?: string;
            created_by?: string;
            id?: number;
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemResponse"][];
            receipt_code?: string;
            status?: string;
            supplier_name?: string;
            updated_at?: string;
            warehouse_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inbound_dto.UpdateInboundRequest": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemRequest"][];
            status?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_inventory_dto.CreateInventoryRequest": {
            available_qty?: number;
            damaged_qty?: number;
            hold_qty?: number;
            product_id: number;
            reserved_qty?: number;
            updated_by?: string;
            warehouse_id: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse": {
            available_qty?: number;
            created_at?: string;
            damaged_qty?: number;
            hold_qty?: number;
            id?: number;
            product_id?: number;
            reserved_qty?: number;
            updated_at?: string;
            updated_by?: string;
            warehouse_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inventory_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_inventory_dto.UpdateInventoryRequest": {
            available_qty?: number;
            damaged_qty?: number;
            hold_qty?: number;
            reason?: string;
            reserved_qty?: number;
            updated_by?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.CreateOrderRequest": {
            assigned_driver_id?: number;
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderItemRequest"][];
            order_code: string;
            receiver_address: string;
            receiver_district?: string;
            receiver_name: string;
            receiver_phone: string;
            receiver_postal_code?: string;
            receiver_province?: string;
            receiver_ward?: string;
            sender_address: string;
            sender_district?: string;
            sender_name: string;
            sender_phone: string;
            sender_postal_code?: string;
            sender_province?: string;
            sender_ward?: string;
            /** @enum {string} */
            status?: "PENDING" | "RESERVED" | "PICKING" | "PACKING" | "SORTING" | "SHIPPING" | "COMPLETED" | "PICKING_UP";
            warehouse_id: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.OrderItemRequest": {
            product_id?: number;
            product_name?: string;
            quantity: number;
            weight_gram?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.OrderItemResponse": {
            id?: number;
            order_id?: number;
            product_id?: number;
            product_name?: string;
            quantity?: number;
            weight_gram?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse": {
            assigned_driver_id?: number;
            created_at?: string;
            created_by?: string;
            id?: number;
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderItemResponse"][];
            order_code?: string;
            receiver_address?: string;
            receiver_district?: string;
            receiver_name?: string;
            receiver_phone?: string;
            receiver_postal_code?: string;
            receiver_province?: string;
            receiver_ward?: string;
            sender_address?: string;
            sender_district?: string;
            sender_name?: string;
            sender_phone?: string;
            sender_postal_code?: string;
            sender_province?: string;
            sender_ward?: string;
            status?: string;
            updated_at?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_order_dto.UpdateOrderRequest": {
            assigned_driver_id?: number;
            order_code?: string;
            receiver_address?: string;
            receiver_district?: string;
            receiver_name?: string;
            receiver_phone?: string;
            receiver_postal_code?: string;
            receiver_province?: string;
            receiver_ward?: string;
            sender_address?: string;
            sender_district?: string;
            sender_name?: string;
            sender_phone?: string;
            sender_postal_code?: string;
            sender_province?: string;
            sender_ward?: string;
            /** @enum {string} */
            status?: "PENDING" | "RESERVED" | "PICKING" | "PACKING" | "SORTING" | "SHIPPING" | "COMPLETED" | "PICKING_UP";
        };
        "my-web-app_com_smart-logistic-hub_internal_product_dto.CreateProductRequest": {
            category?: string;
            created_by?: string;
            dimensions?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.Dimensions"];
            name: string;
            price?: number;
            sku: string;
            weight_gram?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_product_dto.Dimensions": {
            height?: number;
            length?: number;
            width?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse": {
            category?: string;
            created_at?: string;
            created_by?: string;
            height_cm?: number;
            id?: number;
            length_cm?: number;
            name?: string;
            price?: number;
            sku?: string;
            updated_at?: string;
            weight_gram?: number;
            width_cm?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_product_dto.UpdateProductRequest": {
            category?: string;
            created_by?: string;
            dimensions?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.Dimensions"];
            name?: string;
            price?: number;
            weight_gram?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_profile_dto.CreateProfileRequest": {
            name: string;
            phone: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse": {
            created_at?: string;
            /** @description For frontend compatibility */
            display_name?: string;
            id?: number;
            keycloak_user_id?: string;
            name?: string;
            phone?: string;
            /** @description For frontend compatibility */
            user_sub?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_profile_dto.UpdateProfileRequest": {
            name?: string;
            phone?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_role_dto.PermissionResponse": {
            action?: string;
            resource?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_role_dto.RoleResponse": {
            description?: string;
            id?: number;
            name?: string;
            permissions?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_role_dto.PermissionResponse"][];
        };
        "my-web-app_com_smart-logistic-hub_internal_role_dto.UserRolesAndPermissionsResponse": {
            keycloak_user_id?: string;
            roles?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_role_dto.RoleResponse"][];
        };
        "my-web-app_com_smart-logistic-hub_internal_tracking_dto.CreateTrackingEventRequest": {
            driver_code: string;
            lat?: number;
            lng?: number;
            note?: string;
            order_code: string;
            status_update: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_tracking_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse": {
            driver_code?: string;
            id?: number;
            lat?: number;
            lng?: number;
            note?: string;
            order_code?: string;
            status_update?: string;
            timestamp?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_tracking_dto.UpdateTrackingEventRequest": {
            driver_code?: string;
            lat?: number;
            lng?: number;
            note?: string;
            order_code?: string;
            status_update?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.AssignDriverRequest": {
            driver_code: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.CreateTripRequest": {
            driver_code: string;
            status?: string;
            stops: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopRequest"][];
            total_distance_km?: number;
            trip_code: string;
            vehicle_license_plate?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.Location": {
            lat?: number;
            lng?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse": {
            actual_end_at?: string;
            actual_start_at?: string;
            created_at?: string;
            created_by?: string;
            driver_id?: number;
            estimated_duration_min?: number;
            id?: number;
            status?: string;
            stops?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopResponse"][];
            total_distance_km?: number;
            trip_code?: string;
            updated_at?: string;
            vehicle_license_plate?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopRequest": {
            address: string;
            arrived_at?: string;
            departure_at?: string;
            location?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.Location"];
            order_code: string;
            planned_at?: string;
            status?: string;
            stop_type?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopResponse": {
            address?: string;
            arrived_at?: string;
            departure_at?: string;
            id?: number;
            lat?: number;
            lng?: number;
            order_code?: string;
            planned_at?: string;
            status?: string;
            stop_type?: string;
            trip_id?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_trip_dto.UpdateTripRequest": {
            completed_at?: string;
            started_at?: string;
            status?: string;
            stops?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopRequest"][];
        };
        "my-web-app_com_smart-logistic-hub_internal_user_dto.CreateUserRequest": {
            email?: string;
            full_name?: string;
            is_active?: boolean;
            keycloak_sub?: string;
            phone?: string;
            role?: string;
            username: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_user_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_user_dto.UpdateUserRequest": {
            email?: string;
            full_name?: string;
            is_active?: boolean;
            keycloak_sub?: string;
            phone?: string;
            role?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse": {
            created_at?: string;
            created_by?: string;
            email?: string;
            full_name?: string;
            id?: number;
            is_active?: boolean;
            keycloak_sub?: string;
            phone?: string;
            role?: string;
            updated_at?: string;
            username?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_warehouse_dto.CreateWarehouseRequest": {
            address: string;
            contact_phone?: string;
            location?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.Location"];
            manager_name?: string;
            name: string;
            warehouse_code: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_warehouse_dto.Location": {
            lat?: number;
            lng?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_warehouse_dto.UpdateWarehouseRequest": {
            address?: string;
            contact_phone?: string;
            is_active?: boolean;
            location?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.Location"];
            manager_name?: string;
            name?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse": {
            address?: string;
            contact_phone?: string;
            created_at?: string;
            id?: number;
            is_active?: boolean;
            lat?: number;
            lng?: number;
            manager_name?: string;
            name?: string;
            updated_at?: string;
            warehouse_code?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_workspace_dto.CreateWorkspaceRequest": {
            description?: string;
            name: string;
            workspace_code: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_workspace_dto.PaginatedResponse": {
            items?: components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.WorkspaceResponse"][];
            limit?: number;
            skip?: number;
            total?: number;
        };
        "my-web-app_com_smart-logistic-hub_internal_workspace_dto.UpdateWorkspaceRequest": {
            description?: string;
            is_active?: boolean;
            name?: string;
        };
        "my-web-app_com_smart-logistic-hub_internal_workspace_dto.WorkspaceResponse": {
            created_at?: string;
            created_by?: string;
            description?: string;
            id?: number;
            is_active?: boolean;
            name?: string;
            updated_at?: string;
            workspace_code?: string;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export type operations = Record<string, never>;
