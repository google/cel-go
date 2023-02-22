import { TestBed } from '@angular/core/testing';

import { ReplApiService } from './repl-api.service';

describe('ReplApiService', () => {
  let service: ReplApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ReplApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
