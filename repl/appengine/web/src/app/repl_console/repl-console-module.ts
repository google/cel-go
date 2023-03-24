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
