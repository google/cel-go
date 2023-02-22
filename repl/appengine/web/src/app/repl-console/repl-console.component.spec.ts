import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReplConsoleComponent } from './repl-console.component';

describe('ReplConsoleComponent', () => {
  let component: ReplConsoleComponent;
  let fixture: ComponentFixture<ReplConsoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReplConsoleComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReplConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
