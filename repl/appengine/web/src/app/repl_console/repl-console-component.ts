import { Component } from '@angular/core';
import { ReplApiService, EvaluateResponse, EvaluateRequest } from '../shared/repl-api-service';

/**
 * Component for the repl console.
 * Handles input for requests against the REPL api.
 */
@Component({
  selector: 'app-repl-console',
  templateUrl: './repl-console-component.html',
  styleUrls: ['./repl-console-component.scss']
})
export class ReplConsoleComponent {
  lastEvaluate: EvaluateResponse = {responses: [], evalTime: 0};
  lastRequest: EvaluateRequest = {commands: []};

  constructor (private readonly replService: ReplApiService) {}

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

  onEnter(event : KeyboardEvent) : void {
    if (event.key !== "Enter" || event.ctrlKey || event.metaKey) {
      return;
    }
    event.stopPropagation();
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
