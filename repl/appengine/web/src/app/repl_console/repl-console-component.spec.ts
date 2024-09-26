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

import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { ReplConsoleComponent } from './repl-console-component';
import { ReplResultDetailComponent } from './repl-result-detail-component';
import { SharedModule } from '../shared/shared-module';
import { EvaluateRequest } from '../shared/repl-api-service';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

describe('ReplConsoleComponent', () => {
  let component: ReplConsoleComponent;
  let fixture: ComponentFixture<ReplConsoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ReplConsoleComponent, ReplResultDetailComponent],
    imports: [MatFormFieldModule, MatIconModule,
        MatInputModule, SharedModule, NoopAnimationsModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
})
    .compileComponents();

    fixture = TestBed.createComponent(ReplConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
  
  it('default empty', () => {
    expect(component.lastRequest).toEqual(<EvaluateRequest>{commands: []});
  });
});
