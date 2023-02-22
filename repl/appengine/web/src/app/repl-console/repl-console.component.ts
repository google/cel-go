import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { ReplApiService, EvaluateResponse, EvaluateRequest } from '../repl-api.service';

@Component({
  selector: 'app-repl-console',
  templateUrl: './repl-console.component.html',
  styleUrls: ['./repl-console.component.scss']
})
export class ReplConsoleComponent {
  lastEvaluate: EvaluateResponse = {responses: [], evalTime: 0}
  lastRequest: EvaluateRequest = {commands: []}

  constructor (private replService: ReplApiService) {}

  private evaluate(request : EvaluateRequest) {
    this.replService.evaluate(request)
      .subscribe({
        next: (resp) => {
        this.lastRequest = request
        this.lastEvaluate = resp
        let input = document.querySelector<HTMLInputElement>(".repl-stmt-new")
        if (input) { input.value = ""; input.focus() }
        },
      error: (err) => console.log("error: ", err)
      });
  }

  onEnter(event : KeyboardEvent) : void {
    if (event.key !== "Enter" || event.ctrlKey || event.metaKey) {
      return;
    }
    event.stopPropagation()
    let commands : string[] = [];

    document.querySelectorAll(".repl-stmt-input").forEach(
      (el : Element) => {
        if (!(el instanceof HTMLInputElement)) {
          return;
        }
        let inp = el as HTMLInputElement
        if (inp.value && inp.value.trim()) {
          commands.push(inp.value)
        }
      }
    )
    this.evaluate({
      commands: commands
    })
  }

  numStatements() : number {
    return this.lastRequest.commands.length
  }

  focusIndex(index : number) : void  {
    let maxIdx = this.numStatements()
    if (index < 0) {
      index = maxIdx
    } else if (index > maxIdx) {
      index = 0
    }
    document.querySelector<HTMLElement>(`input.repl-stmt-input[data-stmt-index="${index}"]`)?.focus()
  }

  handleUp(i : number) : void {
    this.focusIndex(i - 1)
  }

  handleDown(i : number) : void {
    if (i >= 0) {
      this.focusIndex(i + 1)
    } else {
      this.focusIndex(-1)
    }
  }
}
