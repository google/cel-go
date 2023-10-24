import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { EvaluateRequest } from './repl-api-service';

/**
 * Representation of an example REPL session.
 */
export declare interface Example {
  description: string;
  request: EvaluateRequest;
}

/**
 * Service for setting an active example.
 */
@Injectable({
  providedIn: 'root'
})
export class ReplExampleService {

  private exampleSource = new Subject<Example>();
  
  /**
   * Observable for posted examples.
   */
  examplePosted$ = this.exampleSource.asObservable();

  /**
   * Post an example to be displayed.
   */
  postExample(example: Example) {
    this.exampleSource.next(example);
  }
}
