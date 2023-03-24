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

import { Component, Input } from '@angular/core';
import { CommandResponse } from '../shared/repl-api-service';

/**
 * Simple component for detailing the output of the last command from the REPL
 * API.
 */
@Component({
  selector: 'app-repl-result-detail',
  templateUrl: './repl-result-detail-component.html',
  styleUrls: ['./repl-result-detail-component.scss']
})
export class ReplResultDetailComponent {
  @Input() lastResponse? : CommandResponse;
  @Input() evalTime? : number;

  formatTime(ns : number) : string {
    const ms = ns / 1000000;

    return ms.toPrecision(3) + "ms";
  }
}