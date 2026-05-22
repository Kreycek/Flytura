import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OutPutInvoicesComponent } from './out-put-invoices.component';

describe('OutPutInvoicesComponent', () => {
  let component: OutPutInvoicesComponent;
  let fixture: ComponentFixture<OutPutInvoicesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OutPutInvoicesComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OutPutInvoicesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
