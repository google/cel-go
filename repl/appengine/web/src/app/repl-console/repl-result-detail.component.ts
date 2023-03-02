import { Component, Input } from '@angular/core';
import { CommandResponse } from '../shared/repl-api.service';

/**
 * Simple component for detailing the output of the last command from the REPL
 * API.
 */
@Component({
  selector: 'app-repl-result-detail',
  templateUrl: './repl-result-detail.component.html',
  styleUrls: ['./repl-result-detail.component.scss']
})
export class ReplResultDetailComponent {
  @Input() lastResponse? : CommandResponse;
  @Input() evalTime? : number;

  FormatTime(ns : number) : string {
    const ms = ns / 1000000;

    return ms.toPrecision(3) + "ms";
  }
}