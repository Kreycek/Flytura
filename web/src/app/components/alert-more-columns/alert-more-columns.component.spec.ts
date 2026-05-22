import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AlertMoreColumnsComponent } from './alert-more-columns.component';

describe('AlertMoreColumnsComponent', () => {
  let component: AlertMoreColumnsComponent;
  let fixture: ComponentFixture<AlertMoreColumnsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AlertMoreColumnsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AlertMoreColumnsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
