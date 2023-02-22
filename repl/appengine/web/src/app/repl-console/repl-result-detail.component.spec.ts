import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReplResultDetailComponent } from './repl-result-detail.component';

describe('ReplResultDetailComponent', () => {
  let component: ReplResultDetailComponent;
  let fixture: ComponentFixture<ReplResultDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReplResultDetailComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReplResultDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
