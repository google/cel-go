import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReferencePanelComponent } from './reference-panel.component';

describe('ReferencePanelComponent', () => {
  let component: ReferencePanelComponent;
  let fixture: ComponentFixture<ReferencePanelComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReferencePanelComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReferencePanelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
