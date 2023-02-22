import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient, HttpHeaders} from '@angular/common/http'


const API = "/api";
const HTTP_OPTIONS = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};

export interface EvaluateRequest {
  commands: string[]
}

export interface CommandResponse {
  replOutput: string
  issue: string
  evaluated: boolean
}

export interface EvaluateResponse {
  responses: CommandResponse[]
  evalTime: number
}

@Injectable({
  providedIn: 'root'
})
export class ReplApiService {

  constructor(private httpClient : HttpClient) { }


  evalInternal(request : EvaluateRequest) : Observable<EvaluateResponse> {
    return this.httpClient.post<EvaluateResponse>(API, request, HTTP_OPTIONS)
  }

  evaluate(request : EvaluateRequest) : Observable<EvaluateResponse> {
    return this.evalInternal(request);
  } 
}