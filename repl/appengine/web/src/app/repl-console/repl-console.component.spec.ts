import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { ReplConsoleComponent } from './repl-console.component';
import { ReplResultDetailComponent } from './repl-result-detail.component';
import { SharedModule } from '../shared/shared.module';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { EvaluateRequest } from '../shared/repl-api.service';

describe('ReplConsoleComponent', () => {
  let component: ReplConsoleComponent;
  let fixture: ComponentFixture<ReplConsoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ HttpClientTestingModule, MatFormFieldModule, MatIconModule,
        MatInputModule, SharedModule, NoopAnimationsModule ],
      declarations: [ ReplConsoleComponent, ReplResultDetailComponent ]
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
