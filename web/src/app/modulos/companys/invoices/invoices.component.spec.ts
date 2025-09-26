import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoicesOnlyFlyComponent } from './invoices.component';

describe('InvoicesComponent', () => {
  let component: InvoicesOnlyFlyComponent;
  let fixture: ComponentFixture<InvoicesOnlyFlyComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InvoicesOnlyFlyComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(InvoicesOnlyFlyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
