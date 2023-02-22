import { Component, Input } from '@angular/core';
import { CommandResponse } from '../repl-api.service';

@Component({
  selector: 'app-repl-result-detail',
  templateUrl: './repl-result-detail.component.html',
  styleUrls: ['./repl-result-detail.component.scss']
})
export class ReplResultDetailComponent {
  @Input() lastResponse? : CommandResponse;
  constructor () {}

}