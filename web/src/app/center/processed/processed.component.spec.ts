import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProcessedComponent } from './processed.component';

describe('ProcessedComponent', () => {
  let component: ProcessedComponent;
  let fixture: ComponentFixture<ProcessedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProcessedComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProcessedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
