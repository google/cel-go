import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';


const API = "/api";
const HTTP_OPTIONS = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};

/**
 * Evaluate request for the REPL. A list of repl commands starting with a
 * default instance.
 */
export declare interface EvaluateRequest {
  commands: string[];
}

/**
 * Result of a single REPL command/
 */
export declare interface CommandResponse {
  replOutput: string;
  issue: string;
  evaluated: boolean;
}

/**
 * Result of an evaluate request from the REPL.
 */
export declare interface EvaluateResponse {
  responses: CommandResponse[];
  evalTime: number;
}

/**
 * Wrapper for the CEL REPL JSON API.
 */
@Injectable({
  providedIn: 'root'
})
export class ReplApiService {

  constructor(private readonly httpClient : HttpClient) { }


  evalInternal(request : EvaluateRequest) : Observable<EvaluateResponse> {
    return this.httpClient.post<EvaluateResponse>(API, request, HTTP_OPTIONS);
  }

  evaluate(request : EvaluateRequest) : Observable<EvaluateResponse> {
    return this.evalInternal(request);
  } 
}