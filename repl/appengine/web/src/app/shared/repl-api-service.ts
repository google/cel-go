/**
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { HttpClient, HttpHeaders } from '@angular/common/http';
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