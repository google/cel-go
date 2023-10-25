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

import { Component, OnInit } from '@angular/core';
import { ReplApiService, EvaluateResponse, EvaluateRequest, CommandResponse } from '../shared/repl-api-service';
import { Example, ReplExampleService } from '../shared/repl-example-service';

/**
 * Component for the repl console.
 * Handles input for requests against the REPL api.
 */
@Component({
  selector: 'app-repl-console',
  templateUrl: './repl-console-component.html',
  styleUrls: ['./repl-console-component.scss']
})
export class ReplConsoleComponent implements OnInit {
  lastEvaluate: EvaluateResponse = {responses: [], evalTime: 0};
  lastRequest: EvaluateRequest = {commands: []};

  constructor (
    private readonly replService: ReplApiService,
    private readonly exampleService: ReplExampleService) {}

  ngOnInit() {
    this.exampleService.examplePosted$.subscribe({
      next: (ex: Example) => {
        this.lastRequest = ex.request;
        this.lastEvaluate = {responses: [], evalTime: 0};
        const input = document.querySelector<HTMLInputElement>(".repl-stmt-new");
        if (input) { input.value = ""; input.focus(); }
      }
    });
  }

  private evaluate(request : EvaluateRequest) {
    this.replService.evaluate(request)
      .subscribe({
        next: (resp : EvaluateResponse) => {
        this.lastRequest = request;
        this.lastEvaluate = resp;
        const input = document.querySelector<HTMLInputElement>(".repl-stmt-new");
        if (input) { input.value = ""; input.focus(); }
        },
      error: (err) => { console.log("error: ", err); }
      });
  }

  submit() {
    const request : EvaluateRequest = {commands: []};
    document.querySelectorAll(".repl-stmt-input").forEach(
      (el : Element) => {
        if (!(el instanceof HTMLInputElement)) {
          return;
        }
        if (el.value && el.value.trim()) {
          request.commands.push(el.value);
        }
      }
    );
    this.evaluate(request);
  }

  onEnter(event : KeyboardEvent) : void {
    if (event.key !== "Enter" || event.ctrlKey || event.metaKey) {
      return;
    }
    event.stopPropagation();
    this.submit();
  }

  getResponse(i: number) : CommandResponse {
    if (i < this.lastEvaluate.responses.length) {
      return this.lastEvaluate.responses[i];
    }
    return {replOutput: "", issue: "", evaluated: false};
  }

  reset() : void {
    this.lastEvaluate = {responses: [], evalTime: 0};
    this.lastRequest = {commands: []};
  }

  numStatements() : number {
    return this.lastRequest.commands.length;
  }

  focusIndex(index : number) : void  {
    const maxIdx = this.numStatements();
    if (index < 0) {
      index = maxIdx;
    } else if (index > maxIdx) {
      index = 0;
    }
    document.querySelector<HTMLElement>(`input.repl-stmt-input[data-stmt-index="${index}"]`)?.focus();
  }

  handleUp(i : number) : void {
    this.focusIndex(i - 1);
  }

  handleDown(i : number) : void {
    if (i >= 0) {
      this.focusIndex(i + 1);
    } else {
      this.focusIndex(-1);
    }
  }
}
