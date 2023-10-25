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

import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { EvaluateRequest } from './repl-api-service';

/**
 * Representation of an example REPL session.
 */
export declare interface Example {
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
