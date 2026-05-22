import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PurcharseRecordComponent } from './purcharse-record.component';

describe('PurcharseComponent', () => {
  let component: PurcharseRecordComponent;
  let fixture: ComponentFixture<PurcharseRecordComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PurcharseRecordComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PurcharseRecordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
