/**
 * Wasp API
 * REST API for the Wasp node
 *
 * OpenAPI spec version: 0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { HttpFile } from '../http/http';

export class ContractInfoResponse {
    /**
    * The id (HName as Hex)) of the contract.
    */
    'hName': string;
    /**
    * The name of the contract.
    */
    'name': string;
    /**
    * The hash of the contract. (Hex encoded)
    */
    'programHash': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "hName",
            "baseName": "hName",
            "type": "string",
            "format": "string"
        },
        {
            "name": "name",
            "baseName": "name",
            "type": "string",
            "format": "string"
        },
        {
            "name": "programHash",
            "baseName": "programHash",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return ContractInfoResponse.attributeTypeMap;
    }

    public constructor() {
    }
}
