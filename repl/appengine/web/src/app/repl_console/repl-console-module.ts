import { NgModule } from '@angular/core';
import { ReplConsoleComponent } from './repl-console-component';
import { HttpClientModule } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { ReplResultDetailComponent } from './repl-result-detail-component';
import { SharedModule } from '../shared/shared-module';


@NgModule({
  declarations: [
    ReplConsoleComponent,
    ReplResultDetailComponent,
  ],
  imports: [
    CommonModule,
    HttpClientModule,
    MatFormFieldModule,
    SharedModule,
    MatInputModule,
    MatIconModule,
    MatButtonModule,
  ],
  exports: [
    ReplConsoleComponent
  ]
})
export class ReplConsoleModule { }
