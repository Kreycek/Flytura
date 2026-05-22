import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConciliationListComponent } from './list.component';

describe('ListComponent', () => {
  let component: ConciliationListComponent;
  let fixture: ComponentFixture<ConciliationListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConciliationListComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ConciliationListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
