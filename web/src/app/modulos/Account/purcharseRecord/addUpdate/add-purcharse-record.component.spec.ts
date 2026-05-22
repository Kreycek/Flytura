import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddPurchaseRecordComponent } from './add-purcharse-record.component';

describe('AddInvoicesComponent', () => {
  let component: AddPurchaseRecordComponent;
  let fixture: ComponentFixture<AddPurchaseRecordComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddPurchaseRecordComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AddPurchaseRecordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
